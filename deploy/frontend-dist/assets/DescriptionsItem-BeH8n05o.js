import{Ci as e,Mr as t,Si as n,Ti as r,Zi as i,Zr as a,aa as o,bi as s,br as c,ia as l,jr as u,rr as d,si as f,vi as p,wi as m,xi as h,yi as g,yr as _}from"./router-CIFXj4id.js";import{a as v}from"./index-BSgXDaa4.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=p([g(`descriptions`,{fontSize:`var(--n-font-size)`},[g(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),g(`descriptions-table-wrapper`,[g(`descriptions-table`,[g(`descriptions-table-row`,[g(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),g(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),n(`bordered`,[g(`descriptions-table-wrapper`,[g(`descriptions-table`,[g(`descriptions-table-row`,[p(`&:last-child`,[g(`descriptions-table-content`,{paddingBottom:0})])])])])]),h(`left-label-placement`,[g(`descriptions-table-content`,[p(`> *`,{verticalAlign:`top`})])]),h(`left-label-align`,[p(`th`,{textAlign:`left`})]),h(`center-label-align`,[p(`th`,{textAlign:`center`})]),h(`right-label-align`,[p(`th`,{textAlign:`right`})]),h(`bordered`,[g(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[g(`descriptions-table`,[g(`descriptions-table-row`,[p(`&:not(:last-child)`,[g(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),g(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),g(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[p(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),g(`descriptions-table-content`,[p(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),g(`descriptions-header`,`
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
 `),m(g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),r(g(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=l({name:`Descriptions`,props:Object.assign(Object.assign({},d.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(t){let{mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedComponentPropsRef:o}=c(t),s=i(()=>t.size||o?.value?.Descriptions?.size||`medium`),l=d(`Descriptions`,`-descriptions`,b,v,t,n),u=i(()=>{let{bordered:n}=t,r=s.value,{common:{cubicBezierEaseInOut:i},self:{titleTextColor:a,thColor:o,thColorModal:c,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[e(`fontSize`,r)]:S,[e(n?`thPaddingBordered`:`thPadding`,r)]:C,[e(n?`tdPaddingBordered`:`tdPadding`,r)]:w}}=l.value;return{"--n-title-text-color":a,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":i,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":o,"--n-th-color-modal":c,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),f=r?_(`descriptions`,i(()=>{let e=``,{bordered:n}=t;return n&&(e+=`a`),e+=s.value[0],e}),u,t):void 0;return{mergedClsPrefix:n,cssVars:r?void 0:u,themeClass:f?.themeClass,onRender:f?.onRender,compitableColumn:a(t,[`columns`,`column`]),inlineThemeDisabled:r,mergedSize:s}},render(){let e=this.$slots.default,n=e?t(e()):[];n.length;let{contentClass:r,labelClass:i,compitableColumn:a,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:d,title:p,cssVars:m,mergedClsPrefix:h,separator:g,onRender:_}=this;_?.();let v=n.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let c=t.props||{},l=v.length-1===n,u=[`label`in c?c.label:y(t,`label`)],f=[y(t)],p=c.span||1,m=e.span;e.span+=p;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(o(`th`,{class:[`${h}-descriptions-table-header`,i],colspan:1,style:_},u),o(`td`,{class:[`${h}-descriptions-table-content`,r],colspan:l?(a-m)*2+1:p*2-1,style:b},f)):e.row.push(o(`td`,{class:`${h}-descriptions-table-content`,colspan:l?(a-m)*2:p*2},o(`span`,{class:[`${h}-descriptions-table-content__label`,i],style:_},[...u,g&&o(`span`,{class:`${h}-descriptions-separator`},g)]),o(`span`,{class:[`${h}-descriptions-table-content__content`,r],style:b},f)));else{let t=l?(a-m)*2:p*2;e.row.push(o(`th`,{class:[`${h}-descriptions-table-header`,i],colspan:t,style:_},u)),e.secondRow.push(o(`td`,{class:[`${h}-descriptions-table-content`,r],colspan:t,style:b},f))}return(e.span>=a||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>o(`tr`,{class:`${h}-descriptions-table-row`},e));return o(`div`,{style:m,class:[`${h}-descriptions`,this.themeClass,`${h}-descriptions--${s}-label-placement`,`${h}-descriptions--${c}-label-align`,`${h}-descriptions--${l}-size`,d&&`${h}-descriptions--bordered`]},p||this.$slots.header?o(`div`,{class:`${h}-descriptions-header`},p||u(this,`header`)):null,o(`div`,{class:`${h}-descriptions-table-wrapper`},o(`table`,{class:`${h}-descriptions-table`},o(`tbody`,null,s===`top`&&o(`tr`,{class:`${h}-descriptions-table-row`,style:{visibility:`collapse`}},f(a*2,o(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=l({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};