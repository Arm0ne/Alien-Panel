import{Ai as e,Ci as t,Di as n,Ei as r,Ir as i,Lr as a,Oi as o,Ti as s,Tr as c,cr as l,fi as u,ki as d,la as f,na as p,ni as m,ua as h,wi as g,wr as _}from"./router-BPt_aGUh.js";import{a as v}from"./index-BmdhIupJ.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=t([g(`descriptions`,{fontSize:`var(--n-font-size)`},[g(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),g(`descriptions-table-wrapper`,[g(`descriptions-table`,[g(`descriptions-table-row`,[g(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),g(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),n(`bordered`,[g(`descriptions-table-wrapper`,[g(`descriptions-table`,[g(`descriptions-table-row`,[t(`&:last-child`,[g(`descriptions-table-content`,{paddingBottom:0})])])])])]),r(`left-label-placement`,[g(`descriptions-table-content`,[t(`> *`,{verticalAlign:`top`})])]),r(`left-label-align`,[t(`th`,{textAlign:`left`})]),r(`center-label-align`,[t(`th`,{textAlign:`center`})]),r(`right-label-align`,[t(`th`,{textAlign:`right`})]),r(`bordered`,[g(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[g(`descriptions-table`,[g(`descriptions-table-row`,[t(`&:not(:last-child)`,[g(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),g(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),g(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[t(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),g(`descriptions-table-content`,[t(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),g(`descriptions-header`,`
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
 `,[s(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),s(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),d(g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),e(g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=f({name:`Descriptions`,props:Object.assign(Object.assign({},l.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n,mergedComponentPropsRef:r}=c(e),i=p(()=>e.size||r?.value?.Descriptions?.size||`medium`),a=l(`Descriptions`,`-descriptions`,b,v,e,t),s=p(()=>{let{bordered:t}=e,n=i.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:s,thColor:c,thColorModal:l,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[o(`fontSize`,n)]:S,[o(t?`thPaddingBordered`:`thPadding`,n)]:C,[o(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=a.value;return{"--n-title-text-color":s,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":c,"--n-th-color-modal":l,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),u=n?_(`descriptions`,p(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=i.value[0],t}),s,e):void 0;return{mergedClsPrefix:t,cssVars:n?void 0:s,themeClass:u?.themeClass,onRender:u?.onRender,compitableColumn:m(e,[`columns`,`column`]),inlineThemeDisabled:n,mergedSize:i}},render(){let e=this.$slots.default,t=e?a(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:o,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:d,title:f,cssVars:p,mergedClsPrefix:m,separator:g,onRender:_}=this;_?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,i)=>{let a=t.props||{},c=v.length-1===i,l=[`label`in a?a.label:y(t,`label`)],u=[y(t)],f=a.span||1,p=e.span;e.span+=f;let _=a.labelStyle||a[`label-style`]||this.labelStyle,b=a.contentStyle||a[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(h(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:1,style:_},l),h(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:c?(o-p)*2+1:f*2-1,style:b},u)):e.row.push(h(`td`,{class:`${m}-descriptions-table-content`,colspan:c?(o-p)*2:f*2},h(`span`,{class:[`${m}-descriptions-table-content__label`,r],style:_},[...l,g&&h(`span`,{class:`${m}-descriptions-separator`},g)]),h(`span`,{class:[`${m}-descriptions-table-content__content`,n],style:b},u)));else{let t=c?(o-p)*2:f*2;e.row.push(h(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:t,style:_},l)),e.secondRow.push(h(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:t,style:b},u))}return(e.span>=o||c)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>h(`tr`,{class:`${m}-descriptions-table-row`},e));return h(`div`,{style:p,class:[`${m}-descriptions`,this.themeClass,`${m}-descriptions--${s}-label-placement`,`${m}-descriptions--${c}-label-align`,`${m}-descriptions--${l}-size`,d&&`${m}-descriptions--bordered`]},f||this.$slots.header?h(`div`,{class:`${m}-descriptions-header`},f||i(this,`header`)):null,h(`div`,{class:`${m}-descriptions-table-wrapper`},h(`table`,{class:`${m}-descriptions-table`},h(`tbody`,null,s===`top`&&h(`tr`,{class:`${m}-descriptions-table-row`,style:{visibility:`collapse`}},u(o*2,h(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=f({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};