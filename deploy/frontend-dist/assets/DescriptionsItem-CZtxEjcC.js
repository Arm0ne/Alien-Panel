import{$r as e,Ci as t,Di as n,Ei as r,Nr as i,Pr as a,Si as o,Sr as s,Ti as c,Zi as l,aa as u,ar as d,bi as f,li as p,oa as m,wi as h,xi as g,xr as _}from"./router-XnvCoRTR.js";import{a as v}from"./index-CjMER2w_.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=f([g(`descriptions`,{fontSize:`var(--n-font-size)`},[g(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),g(`descriptions-table-wrapper`,[g(`descriptions-table`,[g(`descriptions-table-row`,[g(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),g(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),h(`bordered`,[g(`descriptions-table-wrapper`,[g(`descriptions-table`,[g(`descriptions-table-row`,[f(`&:last-child`,[g(`descriptions-table-content`,{paddingBottom:0})])])])])]),t(`left-label-placement`,[g(`descriptions-table-content`,[f(`> *`,{verticalAlign:`top`})])]),t(`left-label-align`,[f(`th`,{textAlign:`left`})]),t(`center-label-align`,[f(`th`,{textAlign:`center`})]),t(`right-label-align`,[f(`th`,{textAlign:`right`})]),t(`bordered`,[g(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[g(`descriptions-table`,[g(`descriptions-table-row`,[f(`&:not(:last-child)`,[g(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),g(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),g(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[f(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),g(`descriptions-table-content`,[f(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),g(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),g(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[g(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[g(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[g(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),g(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[o(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),o(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),r(g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),n(g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=u({name:`Descriptions`,props:Object.assign(Object.assign({},d.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(t){let{mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedComponentPropsRef:i}=s(t),a=l(()=>t.size||i?.value?.Descriptions?.size||`medium`),o=d(`Descriptions`,`-descriptions`,b,v,t,n),u=l(()=>{let{bordered:e}=t,n=a.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:s,thColorModal:l,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[c(`fontSize`,n)]:S,[c(e?`thPaddingBordered`:`thPadding`,n)]:C,[c(e?`tdPaddingBordered`:`tdPadding`,n)]:w}}=o.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":s,"--n-th-color-modal":l,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),f=r?_(`descriptions`,l(()=>{let e=``,{bordered:n}=t;return n&&(e+=`a`),e+=a.value[0],e}),u,t):void 0;return{mergedClsPrefix:n,cssVars:r?void 0:u,themeClass:f?.themeClass,onRender:f?.onRender,compitableColumn:e(t,[`columns`,`column`]),inlineThemeDisabled:r,mergedSize:a}},render(){let e=this.$slots.default,t=e?a(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:o,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:u,title:d,cssVars:f,mergedClsPrefix:h,separator:g,onRender:_}=this;_?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,i)=>{let a=t.props||{},c=v.length-1===i,l=[`label`in a?a.label:y(t,`label`)],d=[y(t)],f=a.span||1,p=e.span;e.span+=f;let _=a.labelStyle||a[`label-style`]||this.labelStyle,b=a.contentStyle||a[`content-style`]||this.contentStyle;if(s===`left`)u?e.row.push(m(`th`,{class:[`${h}-descriptions-table-header`,r],colspan:1,style:_},l),m(`td`,{class:[`${h}-descriptions-table-content`,n],colspan:c?(o-p)*2+1:f*2-1,style:b},d)):e.row.push(m(`td`,{class:`${h}-descriptions-table-content`,colspan:c?(o-p)*2:f*2},m(`span`,{class:[`${h}-descriptions-table-content__label`,r],style:_},[...l,g&&m(`span`,{class:`${h}-descriptions-separator`},g)]),m(`span`,{class:[`${h}-descriptions-table-content__content`,n],style:b},d)));else{let t=c?(o-p)*2:f*2;e.row.push(m(`th`,{class:[`${h}-descriptions-table-header`,r],colspan:t,style:_},l)),e.secondRow.push(m(`td`,{class:[`${h}-descriptions-table-content`,n],colspan:t,style:b},d))}return(e.span>=o||c)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>m(`tr`,{class:`${h}-descriptions-table-row`},e));return m(`div`,{style:f,class:[`${h}-descriptions`,this.themeClass,`${h}-descriptions--${s}-label-placement`,`${h}-descriptions--${c}-label-align`,`${h}-descriptions--${l}-size`,u&&`${h}-descriptions--bordered`]},d||this.$slots.header?m(`div`,{class:`${h}-descriptions-header`},d||i(this,`header`)):null,m(`div`,{class:`${h}-descriptions-table-wrapper`},m(`table`,{class:`${h}-descriptions-table`},m(`tbody`,null,s===`top`&&m(`tr`,{class:`${h}-descriptions-table-row`,style:{visibility:`collapse`}},p(o*2,m(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=u({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};